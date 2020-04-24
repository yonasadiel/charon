import React from 'react';
import { SubmissionError } from 'redux-form';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';

import VenueForm, { VenueFormData } from '../../components/exam/venue/form/VenueForm';
import { User, USER_ROLE } from '../../modules/charon/auth/api';
import { Venue } from '../../modules/charon/exam/api';
import * as charonExamActions from '../../modules/charon/exam/action';
import * as charonExamSelectors from '../../modules/charon/exam/selector';
import * as sessionSelectors from '../../modules/session/selector';
import { CharonFormError } from '../../modules/charon/http';
import { AppState } from '../../modules/store';
import { ROUTE_LOGIN } from '../routes';
import './VenueListPage.scss';

interface ConnectedVenueListPageProps {
  venues: Venue[] | null;
  user: User | null;
  getVenues: () => void;
  createVenue: (venue: Venue) => Promise<any>;
};

const renderVenues = (venues: Venue[] | null) => {
  if (venues === null) {
    return (
      <div className="venues">
        <div className="venue"><span className="skeleton">Exam Venue Name</span></div>
        <div className="venue"><span className="skeleton">Another Exam Venue Name</span></div>
        <div className="venue"><span className="skeleton">Even Longer Exam Venue Name</span></div>
      </div>
    );
  }
  if (venues.length === 0) {
    return <div className="venues">Tidak ada lokasi ujian yang sudah ditambahkan.</div>;
  }
  return (
    <div className="venues">
      {venues.map((venue, i) => (
        <div className="venue" key={i}>{venue.name}</div>
      ))}
    </div>
  );
};

const VenueListPage = (props: ConnectedVenueListPageProps) => {
  const { createVenue, venues, getVenues, user } = props;

  React.useEffect(() => { document.title = 'Lokasi Ujian'; }, []);
  React.useEffect(() => { getVenues(); }, [getVenues]);

  const submitNewVenue = async (data: VenueFormData) => {
    return createVenue({ id: 0, ...data } as Venue)
      .then(() => {
        getVenues();
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
    <div className="venue-page">
      <Card>
        <h1 className="title">Daftar Lokasi Ujian</h1>
        {renderVenues(venues)}
        {(user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER) && (
          <div className="add-row">
            <VenueForm onSubmit={submitNewVenue} />
          </div>
        )}
      </Card>

    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  venues: charonExamSelectors.getVenues(state),
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  getVenues: charonExamActions.getVenues,
  createVenue: charonExamActions.createVenue,
};

export default connect(mapStateToProps, mapDispatchToProps)(VenueListPage);
