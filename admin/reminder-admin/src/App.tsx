import {
  Admin,
  Resource,
} from "react-admin";
import { dataProvider } from "./dataProvider";
import authProvider from "./authProvider";

import Reminders from "./components/reminder/index"
import Users from "./components/user/index"
import {Login} from "./LoginPage";

export const App = () => (
  <Admin loginPage={Login} dataProvider={dataProvider} authProvider={authProvider}>
    <Resource
      name="reminders"
      {...Reminders}
    />
    <Resource
        name="users"
        {...Users}
    />
  </Admin>
);
