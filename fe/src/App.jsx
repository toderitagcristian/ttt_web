import { Layout } from "antd";
import { Splash } from "./Splash";
import { useSelector } from "react-redux";
import { CreateOrJoin } from "./CreateOrJoin";
import { Room } from "./Room";

const App = () => {
  const user_token = useSelector((state) => state.user.token);
  const room_id = useSelector((state) => state.room.id)

  return (
    <Layout style={{ height: "100dvh" }}>
      <Layout.Content
        style={{ margin: "16px", border: "1px solid aquamarine" }}
      >
        {!user_token ? <Splash /> : !room_id ? <CreateOrJoin /> : <Room />}
      </Layout.Content>
    </Layout>
  );
};

export default App;
