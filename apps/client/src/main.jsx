import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import Login from './components/Auth/login.jsx';
import SignUp from './components/Auth/SignUp.jsx';
import Room from './pages/Room.jsx';

createRoot(document.getElementById('root')).render(
  <BrowserRouter>
    <Routes>
      <Route path='/' element={<App />} />
      <Route path='login' element={<Login />} />
      <Route path='signup' element={<SignUp />} />
      <Route path='space/:spaceId' element={<Room />} />
      <Route path="*" element={<h1>Not Found</h1>} />
    </Routes>
  </BrowserRouter> 
);
