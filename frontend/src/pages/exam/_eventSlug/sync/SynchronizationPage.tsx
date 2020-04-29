import React from 'react';
import { Card, Button } from 'react-hephaestus';
import { connect } from 'react-redux';
import { withRouter, RouteComponentProps, Redirect } from 'react-router-dom';

import { User } from '../../../../modules/charon/auth/api';
import * as charonExamActions from '../../../../modules/charon/exam/action';
import { SynchronizationData } from '../../../../modules/charon/exam/api';
import * as sessionSelectors from '../../../../modules/session/selector';
import { AppState } from '../../../../modules/store';
import { ROUTE_LOGIN } from '../../../routes';
import './SynchronizationPage.scss';

interface SynchronizationPageProps extends RouteComponentProps<{ eventSlug: string }> {
};

interface ConnectedSynchronizationPageProps extends SynchronizationPageProps {
  user: User | null,
  getSynchronizationData: (eventSlug: string) => Promise<SynchronizationData>,
  putSynchronizationData: (eventSlug: string, syncData: SynchronizationData) => Promise<void>,
};

const SynchronizationPage = (props: ConnectedSynchronizationPageProps) => {
  const {
    getSynchronizationData,
    match: { params: { eventSlug } },
    putSynchronizationData,
    user,
  } = props;

  const [syncData, setSyncData] = React.useState<SynchronizationData>({} as SynchronizationData);
  const handleGetSynchronizationData = () => {
    getSynchronizationData(eventSlug).then((syncData: SynchronizationData) => {
      setSyncData(syncData);
    });
  };
  const handlePutSynchronizationData = React.useCallback(() => {
    putSynchronizationData(eventSlug, syncData).then(() => {
      //
    });
  }, [eventSlug, syncData, putSynchronizationData]);

  if (!user) {
    return <Redirect to={ROUTE_LOGIN} />;
  }

  return (
    <Card className="synchronization-page">
      <h1 className="title">Sinkronisasi</h1>
      <div className="action-row">
        <Button onClick={handleGetSynchronizationData} buttonType="outlined">
          <i className="fas fa-download"/>
          <strong>UNDUH</strong>
        </Button>
        <Button onClick={handlePutSynchronizationData}>
          <i className="fas fa-upload"/>
          <strong>SIMPAN</strong>
        </Button>
      </div>
      <textarea className="sync-data" rows={20} value={JSON.stringify(syncData)} onChange={(e) => setSyncData(JSON.parse(e.currentTarget.value))}/>
    </Card>
  );
};

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  getSynchronizationData: charonExamActions.getSynchronizationData,
  putSynchronizationData: charonExamActions.putSynchronizationData,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(SynchronizationPage));
