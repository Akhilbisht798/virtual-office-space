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
      <Game />
    </>
  )
}

export default App
