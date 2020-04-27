import keyBy from 'lodash/keyBy';
import React from 'react';
import { SubmissionError } from 'redux-form';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Redirect } from 'react-router-dom';

import ParticipationForm, { ParticipationFormData } from '../../../../components/exam/participation/form/ParticipationForm';
import { User, USER_ROLE } from '../../../../modules/charon/auth/api';
import * as charonAuthActions from '../../../../modules/charon/auth/action';
import * as charomAuthSelectors from '../../../../modules/charon/auth/selector';
import * as charonExamActions from '../../../../modules/charon/exam/action';
import { Participation, Venue } from '../../../../modules/charon/exam/api';
import * as charonExamSelectors from '../../../../modules/charon/exam/selector';
import { CharonFormError } from '../../../../modules/charon/http';
import * as sessionSelectors from '../../../../modules/session/selector';
import { AppState } from '../../../../modules/store';
import { ROUTE_LOGIN } from '../../../routes';
import ParticipationLoadingPage from './ParticipationLoadingPage';
import './ParticipationPage.scss';

interface ParticipationPageProps extends RouteComponentProps<{ eventSlug: string }> {
};

interface ConnectedParticipationPageProps extends ParticipationPageProps {
  createParticipation: (eventSlug: string, participation: Participation) => Promise<void>,
  getParticipationsOfEvent: (eventSlug: string) => Promise<void>;
  getUsers: () => Promise<void>,
  getVenues: () => Promise<void>,
  participations: Participation[] | null,
  user: User | null,
  users: User[] | null,
  venues: Venue[] | null,
};

const ParticipationPage = (props: ConnectedParticipationPageProps) => {
  const {
    createParticipation,
    getParticipationsOfEvent,
    getUsers,
    getVenues,
    match: { params: { eventSlug } },
    participations,
    user,
    users,
    venues,
  } = props;

  React.useEffect(() => { if (!participations) getParticipationsOfEvent(eventSlug); }, [getParticipationsOfEvent, eventSlug, participations]);
  React.useEffect(() => { if (!users) getUsers(); }, [getUsers, users]);
  React.useEffect(() => { if (!venues) getVenues(); }, [getVenues, venues]);

  if (!participations) {
    return <ParticipationLoadingPage />;
  }

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }

  const submitNewParticipation = async (data: ParticipationFormData) => {
    return createParticipation(eventSlug, { id: 0, ...data, } as Participation)
      .then(() => {
        getParticipationsOfEvent(eventSlug);
      })
      .catch((err) => {
        if (err instanceof CharonFormError) {
          throw err.asSubmissionError();
        } else {
          throw new SubmissionError({ _error: "Unknown error" });
        }
      });
  };
  const usersByUsername = keyBy(users, 'username');
  const venuesById = keyBy(venues, 'id');

  return (
    <Card className="participation-page">
      <h1 className="title">Daftar Peserta</h1>
      {participations.length === 0 && (<p>Belum ada peserta ditambahkan.</p>)}
      <div className="participations">
        {participations?.map(({userUsername, venueId, id}) => (
          <div className="participation" key={id}>
            <div className="user">{!!usersByUsername[userUsername] ? usersByUsername[userUsername].name : <span className="skeleton">Nama peserta</span>}</div>
            <div className="venue">{!!venuesById[venueId] ? venuesById[venueId].name : <span className="skeleton">Lokasi ujian</span>}</div>
          </div>
        ))}
      </div>
      {(user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER) && (
        <div className="add-row">
          <ParticipationForm
            users={users || []}
            venues={venues || []}
            onSubmit={submitNewParticipation} />
        </div>
      )}
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventSlug: string }>) => ({
  participations: charonExamSelectors.getParticipations(state, props.match.params.eventSlug),
  user: sessionSelectors.getUser(state),
  users: charomAuthSelectors.getUsers(state),
  venues: charonExamSelectors.getVenues(state),
});

const mapDispatchToProps = {
  getParticipationsOfEvent: charonExamActions.getParticipationsOfEvent,
  getUsers: charonAuthActions.getUsers,
  getVenues: charonExamActions.getVenues,
  createParticipation: charonExamActions.createParticipation,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ParticipationPage));
