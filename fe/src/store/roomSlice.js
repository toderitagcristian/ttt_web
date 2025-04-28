import { createSlice } from '@reduxjs/toolkit'

const initialState = {
  id: null,
  p1: null,
  p2: null,
  board: [],
  turn: null,
}

export const roomSlice = createSlice({
  name: 'room',
  initialState,
  reducers: {
    roomCreateSuccess: (state, action) => {
      return {
        ...state,
        ...action.payload,
      }
    },
    joinRoomSuccess: (state, action) => {
      return {
        ...state,
        ...action.payload,
      }
    },
    updateRoomSuccess: (state, action) => {
      return {
        ...state,
        ...action.payload,
      }
    },
  },
})

export const { roomCreateSuccess, joinRoomSuccess, updateRoomSuccess } = roomSlice.actions

export default roomSlice.reducer
