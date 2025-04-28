import { createSlice } from '@reduxjs/toolkit'

const initialState = {
    nickname: "",
    token: null,
}

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    login: (state, action) => {
      return {
          ...state,
          ...action.payload,
      }
    },
    logout: () => ( initialState )
  },
})

export const { login, logout } = userSlice.actions

export default userSlice.reducer