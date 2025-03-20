import { useEffect } from "react"
import { useNavigate } from "react-router"
import { IoEnterOutline } from "react-icons/io5";
import TopBar from "./components/Home/TopBar"

function App() {
  let navigate = useNavigate()

  const handleEnterRoomCode = () => {
    console.log("join room")
    return
  }

  useEffect(() => {
    const jwt = localStorage.getItem("jwt")
    if (!jwt) {
      console.log("jwt not found")
      navigate("/login", { replace: true })
    }
  }, [navigate])

  return (
    <div className="bg-white min-h-screen">
      <TopBar />
      <div className="flex items-center gap-8 p-4">
        <img src="/home/office.webp" alt="office" className="w-1/3 rounded-lg"/>
        <div className="flex flex-col items-start gap-6">
          <h1 className="text-center text-5xl font-semibold text-black">
            Welcome to Zep
          </h1>
          <p className="text-2xl">
            Virtually connect with people in a interactive way
          </p>
        </div>
      </div>
      <div className="gap-4 p-4">
        <div className="flex items-center gap-4">
          <input 
            type="text" 
            placeholder="Enter code"
            className="p-2 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-blue-300" 
          />
          <button 
            className="flex items-center gap-2 bg-blue-500 text-white p-2 rounded"
            onClick={handleEnterRoomCode}>
              <IoEnterOutline/> Join with code
          </button>
          <button 
            className="bg-blue-500 text-white p-2 rounded">
              Create room
          </button>
        </div>
      </div>
    </div>
  )
}

export default App
