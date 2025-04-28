import { useEffect } from 'react';
import useWebSocket from 'react-use-websocket';
import { login, logout } from './store/userSlice';
import { useDispatch } from 'react-redux';
import { roomCreateSuccess, joinRoomSuccess, updateRoomSuccess } from './store/roomSlice';

export const socketUrl = "ws://localhost:3000/ws";
const lsitem = "ttt_token"

const getToken = () => {
    const token = localStorage.getItem(lsitem)
    return token !== null ? token : undefined
}

const setToken = (token) => {
    localStorage.setItem(lsitem, token)
}

const clearToken = () => {
    localStorage.removeItem(lsitem)
}

export const WSClient = () => {
    const dispatch = useDispatch()

    const { 
        lastJsonMessage, 
    } = useWebSocket(socketUrl, {
        shouldReconnect: () => true,
        share: true,
        protocols: getToken()
    });

    useEffect(() => {
        if (lastJsonMessage !== null) {
          switch (lastJsonMessage.event) {
            case "login_success":
                setToken(lastJsonMessage.data.user.token)
                dispatch(login(lastJsonMessage.data.user))
                break;
            case "login_error":
                clearToken()
                dispatch(logout())
                break;
            case "room_create_success":
                dispatch(roomCreateSuccess(lastJsonMessage.data.room))
                break;
            case "room_join_success":
                dispatch(joinRoomSuccess(lastJsonMessage.data.room))
                break;
            case "room_update":
                dispatch(updateRoomSuccess(lastJsonMessage.data.room))
                break;
            default:
                console.log('unknown event: ', lastJsonMessage.event)
          }
        }
    }, [lastJsonMessage, dispatch]);

    return null;
}