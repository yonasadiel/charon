import flatMap from 'lodash/flatMap';
import sortBy from 'lodash/sortBy';

import { AppState } from '../../store';
import { getCharonState } from '../selector';
import { User, USER_ROLE } from './api';
import { CharonAuthState } from './reducer';

function userRoleNumber(user: User): number {
  if (user.role === USER_ROLE.ADMIN) {
    return 10;
  } else if (user.role === USER_ROLE.ORGANIZER) {
    return 20;
  } else if (user.role === USER_ROLE.LOCAL) {
    return 30;
  } else /* user.role === USER_ROLE.PARTICIPANT */ {
    return 40;
  }
}

export function getCharonAuthState(state: AppState): CharonAuthState {
  return getCharonState(state).auth;
}

export function getUsers(state: AppState): User[] | null {
  const { users } = getCharonAuthState(state);
  return !!users ? sortBy(flatMap(users), userRoleNumber) : users;
}
