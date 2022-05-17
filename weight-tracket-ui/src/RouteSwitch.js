import {
  BrowserRouter,
  Routes,
  Route
} from "react-router-dom";

import Users from "./Users";
import User from "./User";

function RouteSwitch(){
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Users />} />
        <Route path="/user/:userId" element={<User />} />
      </Routes>
    </BrowserRouter>
  )
}

export default RouteSwitch;