import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { SERVER } from '../../global';
import { useNavigate } from 'react-router';
import { FaFileUpload } from "react-icons/fa";

const Spaces = () => {
    const [spacesId, setSpacesId] = useState("");
    const [maps, setMaps] = useState([]);
    const [mapId, setMapId] = useState("")
    const [isPublic, setisPublic] = useState(true)
    const [error, setError] = useState("")
    const navigate = useNavigate()

    const handleCreateRoomRequest = async () => {
        if (mapId === "") {
            setError("Please select a map")
            return
        }
        if (spacesId === "") {
            setError("Please name the space")
            return
        }

        const res = await axios.post(`${SERVER}/api/v1/createroom`, {
            name: spacesId, 
            mapId: mapId,
            thumbnail: "https://www.google.com",
            public: true, 
            jwt: localStorage.getItem("jwt")
        });
        navigate('/space/' + res.data.roomId)
    };

    const selectMap = (e) => {
        const id = e.target.id
        console.log("map id", id)
        setMapId(id)
    }

    useEffect(() => {
        const fetchMaps = async () => {
            const res = await axios.get(`${SERVER}/api/v1/getMaps`);
            setMaps(res.data.maps)
        }
        fetchMaps();
    }, [])

    return (
        <div className='p-4 bg-white min-h-screen flex flex-col gap-4'>
            <h3 className='text-2xl font-bold p-4'>Create a new space</h3>
            <div className='flex gap-4'>
                <label htmlFor='spaceName' className='text-lg'>Space Name</label>
                <input
                    className='border-2 border-gray-300 rounded'
                    type='text'
                    placeholder='Enter New Space Name'
                    onChange={(e) => setSpacesId(e.target.value)}
                    id='spaceName'
                />
            </div>
            <div>
                <h2>Select a map</h2>
                <div className='flex flex-row flex-wrap gap-4 p-4 border-2 border-gray-300 rounded'>
                    {
                        maps.map((map) => (
                            <div className="flex flex-col align-center justify-center cursor-pointer" id={map.ID} onClick={selectMap} key={map.ID}>
                                <img src="/home/office.webp" className='rounded-lg w-15' id={map.ID} />
                                {map.Name !== null ? map.Name.split('.')[0] : "Map1"}
                            </div>
                        ))
                    }
                </div>
            </div>
            <div>
                <h3 className='inline-flex gap-2 items-center border-2 p-2'><FaFileUpload /> Upload a thumbnail</h3>
            </div>
            <div className='flex items-center gap-2 mt -2'>
                <input 
                    type='checkbox' id='isPublic' 
                    checked={isPublic}
                    onChange={(e => setisPublic(e.target.checked))}
                    className='cursor-pointer'
                />
                <label htmlFor='isPublic' className='cursor-pointer'>
                    Public
                </label>
            </div>
            <div className='text-red-400'>{error}</div>
            <div>
                <button onClick={handleCreateRoomRequest} className='flex items-center gap-2 bg-blue-500 text-white p-2 rounded'>Create New Space</button>
            </div>
        </div>
    )
}

export default Spaces;
