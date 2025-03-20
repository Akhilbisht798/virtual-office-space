import { useNavigate } from "react-router";
import { CiLogout } from "react-icons/ci";

const TopBar = () => {
    let navigate = useNavigate()
    
    const handleLogout = () => {
        localStorage.removeItem("jwt");
        navigate("/login", { replace: true });
    }
    return (
        <div className="flex bg-blue-600 text-white shadow-md p-4 justify-between">
            <h3 className="text-white">Zep</h3>
            <button className="flex gap-2 items-center text-white hover:text-gray-200" onClick={handleLogout}>
                Logout <CiLogout style={{color: 'white'}}/>
            </button>
        </div>
    )
}

export default TopBar;