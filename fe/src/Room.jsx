import { useSelector } from "react-redux";
import { Flex, Typography } from "antd";
import { useCallback } from "react";
import useWebSocket from "react-use-websocket";
import { socketUrl } from "./WSClient";

const styles = {
    "00": { 
        borderRight: "1px solid yellow",
        borderBottom: "1px solid yellow"
    },
    "02": { 
        borderLeft: "1px solid yellow",
        borderBottom: "1px solid yellow"
    },
    "11": { 
        border: "1px solid yellow",
        margin: "-1px"
    },
    "20": { 
        borderRight: "1px solid yellow",
        borderTop: "1px solid yellow"
    },
    "22": { 
        borderLeft: "1px solid yellow",
        borderTop: "1px solid yellow"
    },
}

export const Room = () => {
    const room = useSelector((state) => state.room)
    const user = useSelector((state) => state.user)

    const { 
            sendJsonMessage,
        } = useWebSocket(socketUrl, {
            share: true
        });

    const getTurnNickname = useCallback(() => {
        const turn = room.turn
        if (!turn) return null
        return room[turn]?.user.nickname

    }, [room])

    const handleClick = (i, j) => {
        if (room[room.turn]?.user.nickname !== user.nickname) {
            console.log("not your turn")
            return
        }

        console.log("coords", i, j)
        sendJsonMessage({
            "event": "board_choice",
            "data": {
                "coords": [i, j]
            }
        })
    }

    return (
        <Flex style={{ height: "100%" }} vertical>
            <Flex gap={72}>
                <div>
                <Typography.Title level={2}>
                    Room ID {room.id}
                </Typography.Title>
                </div>
                
                <div>
                <Typography.Title level={2}>
                    {room.p1.user.nickname} X - 0 {room.p2?.user?.nickname ?? "???"}
                </Typography.Title>
                </div>
            </Flex>

            <div>Turn: {getTurnNickname()}</div>

            <div 
                style={{ 
                    display: "grid", 
                    gridTemplateRows: "repeat(3, 1fr)",
                    height: "300px",
                    width: "300px",
                    alignSelf: "center",
                    marginTop: "6rem",
                }}>
                {room.board.map((row, i) => {
                    return (
                        <div 
                            key={i} 
                            style={{ 
                                display: "grid",
                                gridTemplateColumns: "repeat(3, 1fr)",
                            }}>
                            {row.map((col, j) => {
                                return (
                                    <div 
                                        key={j}
                                        style={{
                                            display: "flex",
                                            justifyContent: "center",
                                            alignItems: "center",
                                            ...styles[`${i}${j}`]
                                        }}
                                        onClick={() => handleClick(i, j)}
                                    >
                                        <Typography.Title level={2}>{col}</Typography.Title>
                                    </div>
                                )
                            })}
                        </div>
                    );
                })}
            </div>
        </Flex>
    )
}
