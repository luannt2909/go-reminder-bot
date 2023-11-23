import {
  Admin,
  Resource,
} from "react-admin";
import { dataProvider } from "./dataProvider";
import Tasks from "./components/task/index"
export const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource
      name="tasks"
      {...Tasks}
    />
  </Admin>
);
