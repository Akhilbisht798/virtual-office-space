import Game from "./Game";
import { useParams } from "react-router";

const Room = () => {
  let param = useParams()
  const spaceId = param.spaceId

  return (
    <>
      <Game spaceId={spaceId} />
    </>
  )
}

export default Room;
