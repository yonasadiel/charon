import React from 'react';
import { SubmissionError } from 'redux-form';
import { Button, Card, Modal } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';

import UserForm, { UserFormData } from '../../components/auth/form/UserForm';
import { User, USER_ROLE } from '../../modules/charon/auth/api';
import * as charonAuthActions from '../../modules/charon/auth/action';
import * as charonAuthSelectors from '../../modules/charon/auth/selector';
import * as sessionSelectors from '../../modules/session/selector';
import { CharonFormError } from '../../modules/charon/http';
import { AppState } from '../../modules/store';
import { ROUTE_LOGIN } from '../routes';
import './UserListPage.scss';

interface ConnectedUserListPageProps {
  users: User[] | null;
  user: User | null;
  getUsers: () => void;
  createUser: (user: User) => Promise<any>;
};

const renderUsers = (users: User[] | null) => {
  if (users === null) {
    return (
      <div className="users">
        <div className="user skeleton"><p>Akun 1</p></div>
        <div className="user skeleton"><p>Akun 2</p></div>
        <div className="user skeleton"><p>Akun 3</p></div>
        <div className="user skeleton"><p>Akun 4</p></div>
      </div>
    );
  }
  if (users.length === 0) {
    return <div className="users">Tidak akun yang terdaftar.</div>;
  }
  return (
    <div className="users">
      {users.map((user, i) => (
        <div className={`user ${user.role}`} key={i}>
          <p>{user.name} <span className="username">{user.username}</span></p>
        </div>
      ))}
    </div>
  );
};

const UserListPage = (props: ConnectedUserListPageProps) => {
  const { createUser, users, getUsers, user } = props;

  React.useEffect(() => { document.title = 'Daftar Akun'; }, []);
  React.useEffect(() => { getUsers(); }, [getUsers]);

  const [isShowingCreateModal, setShowingCreateModal] = React.useState(false);
  const submitNewUser = async (data: UserFormData) => {
    return createUser({ id: 0, ...data } as User)
      .then(() => {
        setShowingCreateModal(false);
        getUsers();
      })
      .catch((err) => {
        if (err instanceof CharonFormError) {
          throw err.asSubmissionError();
        } else {
          throw new SubmissionError({ _error: "Unknown error" });
        }
      });
  };

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }

  return (
    <div className="user-page">
      <Modal isShowing={isShowingCreateModal} closeModal={() => { setShowingCreateModal(false); }}>
        <h2 className="create-user-modal-title">Buat Akun</h2>
        <UserForm onSubmit={submitNewUser} />
      </Modal>

      <Card>
        <div className="title-row">
          <h1 className="title">Daftar Akun</h1>
          {(user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER) && (
            <Button onClick={() => setShowingCreateModal(true)}>
              <i className="fas fa-plus"></i> <span>TAMBAH</span>
            </Button>
          )}
        </div>
        {renderUsers(users)}
      </Card>

    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  users: charonAuthSelectors.getUsers(state),
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  getUsers: charonAuthActions.getUsers,
  createUser: charonAuthActions.createUser,
};

export default connect(mapStateToProps, mapDispatchToProps)(UserListPage);
