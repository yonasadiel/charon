import React from 'react';
import ReactDOM from 'react-dom';
import { HephaestusProvider } from 'react-hephaestus';
import { Provider } from 'react-redux';
import { BrowserRouter as Router } from 'react-router-dom';
import { PersistGate } from 'redux-persist/es/integration/react';

import { persistor, store } from './modules/store';
import App from './pages';
import * as serviceWorker from './serviceWorker';
import './styles/index.scss';

const theme = {
  primaryColor: '#0066b5',
  textColor: '#121215',
  foregroundColor: '#fefefe',
  backgroundColor: '#f0f2f9',
  hoverColor: '#f0f0f0',

  fontSize: 14,
  fontSizeLarge: 16,
  fontSizeSmall: 12,
};

ReactDOM.render(
  <Provider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <HephaestusProvider theme={theme}>
        <Router>
          <App />
        </Router>
      </HephaestusProvider>
    </PersistGate>
  </Provider>,
  document.getElementById('root')
);

serviceWorker.register();
