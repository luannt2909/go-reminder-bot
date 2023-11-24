import {
  Admin,
  Resource,
} from "react-admin";
import { dataProvider } from "./dataProvider";
import Reminders from "./components/reminder/index"
export const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource
      name="reminders"
      {...Reminders}
    />
  </Admin>
);
