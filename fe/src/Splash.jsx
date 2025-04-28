import { Button, Flex, Form, Input } from "antd"
import useWebSocket, { ReadyState } from "react-use-websocket";
import { socketUrl } from "./WSClient";

export const Splash = () => {
    const { 
        sendJsonMessage,
        readyState,
    } = useWebSocket(socketUrl, {
        share: true
    });

    const onFinish = (values) => {
        sendJsonMessage({
            "event": "login",
            "data": {
                "nickname": values.nickname
            }
        })
    }

    const onFinishFailed = errorInfo => {
        console.log('Failed:', errorInfo);
    };

    return (
        <Flex 
            justify="center" 
            align="center"
            style={{ height: "100%" }}
        >
            <Form
                layout="vertical"
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                style={{ maxWidth: 200 }}
                size="large"
            >
                <Form.Item 
                    name="nickname"
                    rules={[
                        { required: true, message: 'Please choose a nickname!' },
                        { min: 4, message: 'Minimum 4 characters!'}
                    ]}
                >
                    <Input placeholder="Your nickname" />
                </Form.Item>

                <Form.Item wrapperCol={{span: 8, offset: 8}}>
                    <Button type="primary" htmlType="submit" disabled={readyState !== ReadyState.OPEN}>
                        Login
                    </Button>
                </Form.Item>
            </Form>
        </Flex>
    )
}