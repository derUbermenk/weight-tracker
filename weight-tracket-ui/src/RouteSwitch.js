import {
  BrowserRouter,
  Routes,
  Route
} from "react-router-dom";

import Users from "./Users";
import User from "./User";
import NewUser from "./NewUser";

function RouteSwitch(){
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Users />} />
        <Route path="/user/:userId" element={<User />} />
        <Route path="/user/new" element={<NewUser />} />
      </Routes>
    </BrowserRouter>
  )
}

export default RouteSwitch;