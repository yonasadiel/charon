import { USER_ROLE, User } from '../../../modules/charon/auth/api';
import {
  ROUTE_EVENT_OVERVIEW,
  ROUTE_EVENT_PARTICIPATION,
  ROUTE_EVENT_QUESTION,
  ROUTE_EVENT_QUESTION_DETAIL,
  ROUTE_EVENT_QUESTION_EDIT,
  ROUTE_EVENT_QUESTION_EDIT_CREATE,
  ROUTE_EVENT_SYNC,
  ROUTE_EVENT_DECRYPT,
} from '../../routes';


export const menuByRole = {
  [USER_ROLE.PARTICIPANT]: [
    ROUTE_EVENT_OVERVIEW,
    ROUTE_EVENT_QUESTION,
    ROUTE_EVENT_QUESTION_DETAIL,
  ],
  [USER_ROLE.LOCAL]: [
    ROUTE_EVENT_OVERVIEW,
    ROUTE_EVENT_PARTICIPATION,
    ROUTE_EVENT_QUESTION,
    ROUTE_EVENT_QUESTION_DETAIL,
    ROUTE_EVENT_SYNC,
    ROUTE_EVENT_DECRYPT,
  ],
  [USER_ROLE.ORGANIZER]: [
    ROUTE_EVENT_OVERVIEW,
    ROUTE_EVENT_PARTICIPATION,
    ROUTE_EVENT_QUESTION,
    ROUTE_EVENT_QUESTION_DETAIL,
    ROUTE_EVENT_QUESTION_EDIT,
    ROUTE_EVENT_QUESTION_EDIT_CREATE,
  ],
  [USER_ROLE.ADMIN]: [
    ROUTE_EVENT_OVERVIEW,
    ROUTE_EVENT_PARTICIPATION,
    ROUTE_EVENT_QUESTION,
    ROUTE_EVENT_QUESTION_DETAIL,
    ROUTE_EVENT_QUESTION_EDIT,
    ROUTE_EVENT_QUESTION_EDIT_CREATE,
  ],
}

export function hasPermissionForMenu(user: User | null, route: string): boolean {
  if (!user) return false;
  return menuByRole[user.role].includes(route);
}
