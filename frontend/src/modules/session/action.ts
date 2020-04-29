import { User } from '../charon/auth/api';

export const PUT_USER = 'session/PUT_USER';
export const putUser = (user: User | null) => ({
  type: PUT_USER,
  user,
});

export const PUT_PARTICIPATION_KEY = 'session/PUT_PARTICIPATION_KEY';
export const putParticipationKey = (eventSlug: string, participationKey: string) => ({
  type: PUT_PARTICIPATION_KEY,
  eventSlug,
  participationKey,
});
