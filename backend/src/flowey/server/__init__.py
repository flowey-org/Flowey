from datetime import UTC, datetime
from typing import ClassVar, TypeAliasType
from uuid import UUID, uuid4

from fastapi import FastAPI, HTTPException, Request, Response
from pydantic import BaseModel

server = FastAPI()


class Session(BaseModel):
    ID: ClassVar = TypeAliasType("ID", UUID)

    id: ID = uuid4()
    created_at: datetime = datetime.now(tz=UTC)


Sessions = TypeAliasType("Sessions", dict[Session.ID, Session])
sessions: Sessions = {}


@server.get("/sessions")
async def get_sessions() -> Sessions:
    return sessions


@server.post("/sessions", status_code=201)
async def post_session(request: Request, response: Response) -> Session:
    session = Session()
    sessions[session.id] = session
    response.headers["Location"] = f"{request.url}/{session.id}"
    return session


@server.get("/sessions/{session_id}")
async def get_session(session_id: Session.ID) -> Session:
    try:
        return sessions[session_id]
    except KeyError:
        raise HTTPException(status_code=404) from None


@server.delete(
    "/sessions/{session_id}",
    status_code=204,
    responses={410: {}},
    response_class=Response,
)
async def delete_session(session_id: Session.ID) -> None:
    try:
        sessions.pop(session_id)
    except KeyError:
        raise HTTPException(status_code=410) from None
