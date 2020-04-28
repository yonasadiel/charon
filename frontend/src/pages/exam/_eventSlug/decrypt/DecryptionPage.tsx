import React from 'react';
import { Card, Button, TextInput } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Redirect } from 'react-router-dom';

import { User } from '../../../../modules/charon/auth/api';
import * as charonExamActions from '../../../../modules/charon/exam/action';
import * as sessionSelectors from '../../../../modules/session/selector';
import { AppState } from '../../../../modules/store';
import { ROUTE_LOGIN } from '../../../routes';
import './DecryptionPage.scss';

interface DecryptionPageProps extends RouteComponentProps<{ eventSlug: string }> {
};

interface ConnectedDecryptionPageProps extends DecryptionPageProps {
  user: User | null,
  decryptEvent: (eventSlug: string, key: string) => Promise<void>,
  decryptEventLocal: (eventSlug: string, key: string) => void,
};

const DecryptionPage = (props: ConnectedDecryptionPageProps) => {
  const {
    decryptEvent,
    decryptEventLocal,
    match: { params: { eventSlug } },
    user,
  } = props;

  const [key, setKey] = React.useState('');
  const handleDecryptDatabase = () => {
    decryptEvent(eventSlug, key).then(() => {
      //
    });
  };
  const handleLocalDecryption = () => {
    decryptEventLocal(eventSlug, key);
  };

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }

  return (
    <Card className="synchronization-page">
      <h1 className="title">Dekripsi</h1>
      <div className="action-row">
        <Button onClick={handleDecryptDatabase} buttonType="outlined">
          <i className="fas fa-server"/>
          <strong>DEKRIPSI DB</strong>
        </Button>
        <Button onClick={handleLocalDecryption}>
          <i className="fas fa-download"/>
          <strong>DEKRIPSI LOKAL</strong>
        </Button>
      </div>
      <TextInput
        wrapperStyle={{ width: '100%' }}
        placeholder="Kunci rahasia"
        value={key}
        onChange={(e) => setKey(e.currentTarget.value)} />
    </Card>
  );
};

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  decryptEvent: charonExamActions.decryptEvent,
  decryptEventLocal: charonExamActions.decryptEventLocal,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(DecryptionPage));
