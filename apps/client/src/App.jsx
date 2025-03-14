import Login from "./components/Auth/login"
import Game from "./pages/Game"

function App() {
  const jwt = localStorage.getItem("jwt")
  if (!jwt) {
    console.log("no jwt: ", jwt)
    return <Login />
  }
  return (
    <>
      <Game spaceId="76417741-22c4-42d8-8d0f-3ad0a3f7c27b"/>
    </>
  )
}

export default App
