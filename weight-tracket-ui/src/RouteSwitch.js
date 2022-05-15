import {
  BrowserRouter,
  Routes,
  Route
} from "react-router-dom";

import Users from "./Users";
import Hello from "./Hello";

function RouteSwitch(){
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Users />} />
        <Route path="/hello" element={<Hello />} />
      </Routes>
    </BrowserRouter>
  )
}

export default RouteSwitch;