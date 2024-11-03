<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";

import Button from "@/components/Button.vue";

import CancelButton from "@/components/buttons/CancelButton.vue";

import { state, store } from "@/store";
import { wss } from "@/wss";

const isModalOpen = ref(false);
const isLoggedIn = ref(false);

const isSubmitting = ref(false);
const password = ref("");
const errorMessage = ref("");

const wssStatus = ref(wss.status());

function openModal() {
  isModalOpen.value = true;
}

function closeModal() {
  isModalOpen.value = false;
  password.value = "";
  errorMessage.value = "";
}

function checkLoginStatus() {
  const cookies = document.cookie.split(";");
  for (const cookie of cookies) {
    if (cookie.startsWith("flowey_session_key_present")) {
      password.value = "";
      isLoggedIn.value = true;
      return;
    }
  }
  isLoggedIn.value = false;
}

async function handleSubmit() {
  if (isSubmitting.value) return;
  isSubmitting.value = true;
  errorMessage.value = "";

  let response: Response;
  try {
    response = await fetch(state.endpoint.value + "session/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: state.username.value,
        password: password.value,
      }),
      credentials: "include",
    });
  } catch {
    errorMessage.value = "Login failed: refused to connect to the endpoint.";
    isSubmitting.value = false;
    return;
  }

  if (response.ok) {
    checkLoginStatus();
    if (!isLoggedIn.value) {
      errorMessage.value = "Login failed: no session key received.";
    }
  } else {
    switch (response.status) {
      case 400:
        errorMessage.value = "Login failed: sent a bad request.";
        break;
      case 401:
        errorMessage.value = "Login failed: invalid credentials.";
        break;
      case 404:
        errorMessage.value = "Login failed: couldn't reach the endpoint.";
        break;
      case 500:
        errorMessage.value = "Login failed: internal server error.";
        break;
      default:
        errorMessage.value = "Login failed: unexpected response.";
    }
  }

  isSubmitting.value = false;
}

async function handleLogout() {
  if (isSubmitting.value) return;
  isSubmitting.value = true;
  errorMessage.value = "";

  let response: Response;
  try {
    response = await fetch(state.endpoint.value + "session/", {
      method: "DELETE",
      credentials: "include",
    });
  } catch {
    errorMessage.value = "Logout failed: refused to connect to the endpoint.";
    isSubmitting.value = false;
    return;
  }

  if (response.ok) {
    checkLoginStatus();
    if (isLoggedIn.value) {
      errorMessage.value = "Logout failed: session key is not invalidated.";
    }
  } else {
    switch (response.status) {
      case 404:
        errorMessage.value = "Login failed: couldn't reach the endpoint.";
        break;
      case 500:
        errorMessage.value = "Login failed: internal server error.";
        break;
      default:
        errorMessage.value = "Login failed: unexpected response.";
    }
  }

  isSubmitting.value = false;
}

function checkWsStatus() {
  wssStatus.value = wss.status();
}

onMounted(async () => {
  await store.ready;

  watch(isLoggedIn, (isLoggedIn) => {
    if (isLoggedIn) {
      wss.connect();
    } else {
      wss.disconnect();
    }
  });

  watch(isModalOpen, (isOpen) => {
    if (isOpen) {
      checkLoginStatus();
      checkWsStatus();
    }
  });

  wss.onStatusChange(() => {
    checkWsStatus();
  });

  checkLoginStatus();
  checkWsStatus();
});
</script>

<template>
  <Button
    :class="(!isLoggedIn || wssStatus !== 'connected') && 'suggested'"
    @click="openModal"
  >
    <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
      <path class="background" d="M 22.153846,1.8461538 H 1.8461538 V 22.153846 H 22.153846 Z" fill="none" />
      <!-- eslint-disable-next-line max-len -->
      <path d="M 5.0999581,9.4124843 V 7.6874758 H 15.450022 v -1.725012 h 1.725011 V 7.687476 h 1.725009 V 9.4124843 H 17.175033 V 11.137495 H 15.450022 V 9.4124843 Z M 15.450022,11.137495 h -1.725011 v 1.725011 h 1.725011 z m 0,-5.1750312 h -1.725011 v -1.72501 h 1.725011 z m 3.45002,10.3500612 V 14.587516 H 8.5499797 v -1.72501 H 10.274991 V 11.137495 H 8.5499797 v 1.725011 H 6.8249694 v 1.72501 H 5.0999581 v 1.725009 h 1.7250114 v 1.725013 h 1.7250102 v 1.725007 H 10.274991 V 18.037538 H 8.5499797 v -1.725013 z" fill="black" />
      <path d="M 1.8461538,0 V 1.8461538 H 22.153846 V 0 Z" fill="#000000" />
      <path d="M 24,1.8461538 H 22.153846 V 22.153846 H 24 Z" fill="#000000" />
      <path d="M 1.8461538,22.153846 V 24 H 22.153846 v -1.846154 z" fill="#000000" />
      <path d="M 0,22.153846 H 1.8461538 V 1.8461538 H 0 Z" fill="#000000" />
    </svg>
  </Button>
  <Teleport to="body">
    <div v-if="isModalOpen" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Sync</h2>
          <CancelButton @click="closeModal" />
        </div>
        <template v-if="!isLoggedIn">
          <form @submit.prevent="handleSubmit">
            <div class="form-group">
              <label for="endpoint">Endpoint:</label>
              <input id="endpoint" v-model="state.endpoint.value" type="url" required>
            </div>
            <div class="form-group">
              <label for="username">Username:</label>
              <input id="username" v-model="state.username.value" type="text" required>
            </div>
            <div class="form-group">
              <label for="password">Password:</label>
              <input id="password" v-model="password" type="password" required>
            </div>
            <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
            <button type="submit" :disabled="isSubmitting">Login</button>
          </form>
        </template>
        <template v-else>
          <p>You're logged in as <em>{{ state.username.value }}</em>.</p>
          <p>WSS status: <em>{{ wssStatus }}</em>.</p>
          <form @submit.prevent="handleLogout">
            <div class="form-group">
              <label for="endpoint">Endpoint:</label>
              <input id="endpoint" v-model="state.endpoint.value" type="url" required>
            </div>
            <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
            <button type="submit" :disabled="isSubmitting">Logout</button>
          </form>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<style>
:root {
  --border-size: 3px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal-content {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  max-width: 500px;
  width: calc(90% - 40px);
  max-height: calc(90vh - 40px);
  overflow-y: auto;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  margin: 20px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header > button > svg {
  height: 2em;
  vertical-align: middle;
  padding: 0;
}

.modal-content > form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.modal-content > form > .form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.modal-content > form > .form-group > input {
  --s: var(--border-size);
  box-sizing: border-box;
  padding: 0.5rem;
  border: solid var(--s) #ccc;
  border-radius: 4px;
  height: 2rem;
  clip-path: polygon(
    0 var(--s), var(--s) var(--s), var(--s) 0,
    calc(100% - var(--s)) 0, calc(100% - var(--s)) var(--s), 100% var(--s),
    100% calc(100% - var(--s)), calc(100% - var(--s)) calc(100% - var(--s)), calc(100% - var(--s)) 100%,
    var(--s) 100%, var(--s) calc(100% - var(--s)), 0 calc(100% - var(--s))
  );
}

.modal-content > form > button[type="submit"] {
  --s: var(--border-size);
  box-sizing: border-box;
  background-color: var(--color-button);
  padding: 0.5rem;
  border: solid var(--s);
  border-radius: 4px;
  cursor: pointer;
  height: 2rem;
  clip-path: polygon(
    0 var(--s), var(--s) var(--s), var(--s) 0,
    calc(100% - var(--s)) 0, calc(100% - var(--s)) var(--s), 100% var(--s),
    100% calc(100% - var(--s)), calc(100% - var(--s)) calc(100% - var(--s)), calc(100% - var(--s)) 100%,
    var(--s) 100%, var(--s) calc(100% - var(--s)), 0 calc(100% - var(--s))
  );
}

.modal-content > form > button[type="submit"]:is(:focus-visible, :hover) {
  filter: brightness(0.9);
}

.modal-content > form > button[type="submit"]:disabled {
  background-color: var(--color-button-disabled);
  cursor: not-allowed;
}

.modal-content > form > .error {
  color: red;
}
</style>
