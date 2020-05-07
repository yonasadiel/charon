import keyBy from 'lodash/keyBy';
import React from 'react';
import { Card, Button } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Redirect } from 'react-router-dom';

import { User } from '../../../../modules/charon/auth/api';
import * as charonAuthActions from '../../../../modules/charon/auth/action';
import * as charomAuthSelectors from '../../../../modules/charon/auth/selector';
import * as charonExamActions from '../../../../modules/charon/exam/action';
import { ParticipationStatus } from '../../../../modules/charon/exam/api';
import * as charonExamSelectors from '../../../../modules/charon/exam/selector';
import * as sessionSelectors from '../../../../modules/session/selector';
import { AppState } from '../../../../modules/store';
import { ROUTE_LOGIN } from '../../../routes';
import ParticipationStatusLoadingPage from './ParticipationStatusLoadingPage';
import './ParticipationStatusPage.scss';
import LoadingCircle from '../../../../components/loading/Circle';

interface ParticipationStatusPageProps extends RouteComponentProps<{ eventSlug: string }> {
};

interface ConnectedParticipationStatusPageProps extends ParticipationStatusPageProps {
  lockUser: (username: string) => Promise<void>;
  unlockUser: (username: string) => Promise<void>;
  getParticipationStatus: (eventSlug: string) => Promise<void>;
  deleteParticipationStatus: (eventSlug: string, sessionId: number) => Promise<void>;
  getUsers: () => Promise<void>,
  participationStatus: ParticipationStatus[] | null,
  user: User | null,
  users: User[] | null,
};

const ParticipationStatusPage = (props: ConnectedParticipationStatusPageProps) => {
  const {
    lockUser,
    unlockUser,
    getParticipationStatus,
    deleteParticipationStatus,
    getUsers,
    match: { params: { eventSlug } },
    participationStatus,
    user,
    users,
  } = props;

  React.useEffect(() => { if (!participationStatus) getParticipationStatus(eventSlug); }, [getParticipationStatus, eventSlug, participationStatus]);
  React.useEffect(() => { if (!users) getUsers(); }, [getUsers, users]);
  const [isLocking, setIsLocking] = React.useState(false);

  if (!participationStatus) {
    return <ParticipationStatusLoadingPage />;
  }

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }
  const usersByUsername = keyBy(users, 'username');
  const handleLockUser = (username: string) => {
    return () => {
      setIsLocking(true);
      lockUser(username).then(() => {
        setIsLocking(false);
        getParticipationStatus(eventSlug);
      });
    };
  }
  const handleUnlockUser = (username: string) => {
    return () => {
      setIsLocking(true);
      unlockUser(username).then(() => {
        setIsLocking(false);
        getParticipationStatus(eventSlug);
      });
    };
  };
  const handleDeleteSession = (sessionId: number) => {
    return () => {
      setIsLocking(true);
      deleteParticipationStatus(eventSlug, sessionId).then(() => {
        setIsLocking(false);
        getParticipationStatus(eventSlug);
      });
    };
  };

  return (
    <Card className="participation-status-page">
      <h1 className="title">Status Peserta</h1>
      {participationStatus.length === 0 && (<p>Belum ada peserta ditambahkan.</p>)}
      <div className="participation-status-list">
        <div className="participation-status">
          <div className="user">Nama</div>
          <div className="ip-address">Alamat IP</div>
          <div className="login-at">Waktu Login</div>
          <div className="lock-button">Kunci Sesi</div>
          <div className="delete"></div>
        </div>
        {participationStatus?.map(({ userUsername, ipAddress, loginAt, sessionId, userSessionLocked }) => (
          <div className="participation-status" key={sessionId}>
            <div className="user">{!!usersByUsername[userUsername] ? usersByUsername[userUsername].name : <span className="skeleton">Nama peserta</span>}</div>
            <div className="ip-address">{ipAddress}</div>
            <div className="login-at">{!!loginAt ? loginAt.toLocaleTimeString() : '-'}</div>
            <div className="lock-button">
              {isLocking
                ? <Button><LoadingCircle /></Button>
                : userSessionLocked
                  ? <Button disabled={isLocking} onClick={handleUnlockUser(userUsername)}><i className="fas fa-lock-open"/><span>Unlock</span></Button>
                  : <Button disabled={isLocking} onClick={handleLockUser(userUsername)}><i className="fas fa-lock"/><span>Lock</span></Button>}
            </div>
            <div className="delete"><Button className="delete-button" disabled={isLocking} onClick={handleDeleteSession(sessionId)}><i className="fas fa-trash"/></Button></div>
          </div>
        ))}
      </div>
    </Card>
  );
};

const mapStateToProps = (state: AppState, props: RouteComponentProps<{ eventSlug: string }>) => ({
  participationStatus: charonExamSelectors.getParticipationStatus(state, props.match.params.eventSlug),
  user: sessionSelectors.getUser(state),
  users: charomAuthSelectors.getUsers(state),
});

const mapDispatchToProps = {
  lockUser: charonAuthActions.lockUser,
  unlockUser: charonAuthActions.unlockUser,
  getParticipationStatus: charonExamActions.getParticipationStatus,
  deleteParticipationStatus: charonExamActions.deleteParticipationStatus,
  getUsers: charonAuthActions.getUsers,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(ParticipationStatusPage));
