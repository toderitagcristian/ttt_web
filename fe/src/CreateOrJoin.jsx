import { useSelector } from "react-redux";
import { Button, Flex, Form, Input, Typography } from "antd";
import useWebSocket from "react-use-websocket";
import { socketUrl } from "./WSClient";

export const CreateOrJoin = () => {
  const { sendJsonMessage } = useWebSocket(socketUrl, {
    share: true,
  });

  const nickname = useSelector((state) => state.user.nickname);

  const onFinish = (values) => {
    sendJsonMessage({
        "event": "room_join",
        "data": { id: values.room_id }
    })
  };

  const onFinishFailed = (errorInfo) => {
    console.log("Failed join:", errorInfo);
  };

  const handleCreateRoom = () => {
    sendJsonMessage({
      event: "room_create",
    });
  };

  return (
    <Flex justify="center" align="center" style={{ height: "100%" }} vertical>
      <Typography.Title>Welcome {nickname} !</Typography.Title>

      <Button type="primary" onClick={handleCreateRoom}>
        Create room
      </Button>

      <Typography.Title level={3}>OR</Typography.Title>

      <Form
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        size="large"
        layout="inline"
      >
        <Form.Item
          style={{ maxWidth: 150 }}
          name="room_id"
          rules={[
            { required: true, message: "Please enter room ID!" },
            { len: 4, message: "Exactly 4 characters!" },
          ]}
        >
          <Input placeholder="Room ID" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit">
            Join
          </Button>
        </Form.Item>
      </Form>
    </Flex>
  );
};
