import React from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import { SubmissionError } from 'redux-form';

import * as charonAuthActions from '../../modules/charon/auth/action';
import { CharonFormError } from '../../modules/charon/http';
import * as sessionSelectors from '../../modules/session/selector';
import { AppState } from '../../modules/store';
import { ROUTE_EXAM } from '../routes';
import LoginForm, { LoginFormData } from '../../components/auth/form/LoginForm';
import './LoginPage.scss';

interface LoginPageProps {
  isLoggedIn: boolean;
  loginAction: (username: string, password: string) => Promise<void>,
};

const LoginPage = (props: LoginPageProps) => {
  const { isLoggedIn, loginAction } = props;
  React.useEffect(() => { document.title = 'Login'; }, []);

  const submitLogin = async (data: LoginFormData) => {
    return loginAction(data.username, data.password)
      .then(() => { })
      .catch((err) => {
        if (err instanceof CharonFormError) {
          throw err.asSubmissionError();
        } else {
          throw new SubmissionError({ _error: "Unknown error" });
        }
      });
  };

  if (isLoggedIn) {
    return <Redirect to={ROUTE_EXAM} />;
  }

  return (
    <div className="login-page">
      <Card className="login-card">
        <h1 className="title">Login</h1>
        <LoginForm onSubmit={submitLogin} />
      </Card>
    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  isLoggedIn: sessionSelectors.isLoggedIn(state),
});

const mapDispatchToProps = {
  loginAction: charonAuthActions.login,
};

export default connect(mapStateToProps, mapDispatchToProps)(LoginPage);
