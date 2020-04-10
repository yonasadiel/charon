import { applyMiddleware, combineReducers, createStore, compose, Action } from 'redux';
import { reducer as formReducer, FormStateMap } from 'redux-form';
import { persistStore, persistReducer, PersistState } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import thunk, { ThunkAction } from 'redux-thunk';

import { CharonState, charonReducer } from './charon/reducer';
import charonAuthApi, { CharonAuthApi } from './charon/auth/api';
import charonExamApi, { CharonExamApi } from './charon/exam/api';
import { SessionState, sessionReducer } from './session/reducer';

type PersistedState<T> = T & { _persist: PersistState };

export interface AppState {
  charon: CharonState,
  form: FormStateMap,
  session: PersistedState<SessionState>,
};

export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  AppState,
  ThunkExtraArguments,
  Action<string>
>;

export type ThunkExtraArguments = {
  charonAuthApi: CharonAuthApi,
  charonExamApi: CharonExamApi,
};

const rootReducer = combineReducers<AppState>({
  charon: charonReducer,
  form: formReducer,
  session: persistReducer({ key: 'charonSession', storage }, sessionReducer),
});

const composeEnhancers = (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

export const store: any = createStore<AppState, any, any, any>(
  rootReducer,
  composeEnhancers(
    applyMiddleware(
      thunk.withExtraArgument({
        charonAuthApi,
        charonExamApi,
      }),
    ),
  ),
);

export const persistor = persistStore(store);
