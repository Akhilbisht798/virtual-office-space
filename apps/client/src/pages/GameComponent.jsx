import Game from "./Game";
import { useParams } from "react-router";

const GameComponent = () => {
  let param = useParams()
  const spaceId = param.spaceId

  return (
    <>
      <Game spaceId={spaceId} />
      {/*
      <Game spaceId="76417741-22c4-42d8-8d0f-3ad0a3f7c27b"/>
      */}
    </>
  )
}

export default GameComponent;
