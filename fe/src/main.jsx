import { createRoot } from 'react-dom/client'
import App from './App.jsx'
import "./reset.css"
import { Provider } from 'react-redux'
import { store } from './store/store.js'
import { ConfigProvider, theme } from 'antd'
import { WSClient } from './WSClient.jsx'

createRoot(document.getElementById('root')).render(
    <Provider store={store}>
      <ConfigProvider
        theme={{
          algorithm: theme.darkAlgorithm,
        }}
      >
        <WSClient />
        <App />
      </ConfigProvider>
    </Provider>
)
