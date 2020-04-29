import { User } from '../charon/auth/api';
import { AppState } from '../store';
import { SessionState } from './reducer';

export function getSessionState(state: AppState): SessionState {
  return state.session;
}

export function getUser(state: AppState): User | null {
  return getSessionState(state).user;
}

export function isLoggedIn(state: AppState): boolean {
  return !!getSessionState(state).user;
}

export function getParticipationKey(state: AppState, eventSlug: string): string | undefined {
  return getSessionState(state).participationKey[eventSlug];
}
