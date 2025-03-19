import { useEffect } from "react"
import Game from "./pages/Game"
import { useNavigate } from "react-router"

function App() {
  let navigate = useNavigate()
  useEffect(() => {
    const jwt = localStorage.getItem("jwt")
    if (!jwt) {
      console.log("jwt not found")
      navigate("/login", { replace: true })
    }
  }, [navigate])

  return (
    <>
      <h1>This is home page</h1>
      {/*
      <Game spaceId="76417741-22c4-42d8-8d0f-3ad0a3f7c27b"/>
      */}
    </>
  )
}

export default App
