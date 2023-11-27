import {
  Admin,
  Resource,
} from "react-admin";
import { dataProvider } from "./dataProvider";
import authProvider from "./authProvider";

import Reminders from "./components/reminder/index"
import Users from "./components/user/index"
import {Login} from "./LoginPage";
import ReminderIcon from "@material-ui/icons/NotificationImportant";
import PersonIcon from "@material-ui/icons/Person";
import {RoleAdmin} from "./components/user/role";

export const App = () => (
  <Admin loginPage={Login} dataProvider={dataProvider} authProvider={authProvider}>
      {permissions => (
          <>
              <Resource
                  name="reminders"
                  icon={ReminderIcon}
                  {...Reminders}
              />
              {permissions == RoleAdmin ? <Resource
                  name="users"
                  icon={PersonIcon}
                  {...Users}
              /> : null}

          </>
      )}

  </Admin>
);
