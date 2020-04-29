import React from 'react';
import { SubmissionError } from 'redux-form';
import { Button, Card, Modal } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import EventForm, { EventFormData } from '../../components/exam/event/form/EventForm';
import { User, USER_ROLE } from '../../modules/charon/auth/api';
import { Event } from '../../modules/charon/exam/api';
import * as charonExamActions from '../../modules/charon/exam/action';
import * as charonExamSelectors from '../../modules/charon/exam/selector';
import { CharonFormError } from '../../modules/charon/http';
import { generateUrlWithParams } from '../../modules/util/routes';
import { AppState } from '../../modules/store';
import { ROUTE_EVENT } from '../routes';
import { durationText } from '../../modules/util/time';
import './EventListPage.scss';

interface EventListPageProps {
  events: Event[] | null;
  user: User;
  getEvents: () => void;
  createEvent: (event: Event) => Promise<any>;
};

const renderEvents = (events: Event[] | null) => {
  if (events === null) {
    return (
      <div>
        <Card className="event-card">
          <h2 className="event-title"><span className="skeleton">Exam #1</span></h2>
          <p><span className="skeleton">Starts at 01/01/2020 00:00</span></p>
        </Card>
        <Card className="event-card">
          <h2 className="event-title"><span className="skeleton">Exam #1</span></h2>
          <p><span className="skeleton">Starts at 01/01/2020 00:00</span></p>
        </Card>
        <Card className="event-card">
          <h2 className="event-title"><span className="skeleton">Exam #1</span></h2>
          <p><span className="skeleton">Starts at 01/01/2020 00:00</span></p>
        </Card>
      </div>
    );
  }
  if (events.length === 0) {
    return <p>No Exam</p>;
  }
  return (
    <div>
      {events.map((event: Event) => (
        <Link to={generateUrlWithParams(ROUTE_EVENT, { eventSlug: event.slug })} key={event.id}>
          <Card className="event-card">
            <h2 className="event-title">{event.title}</h2>
            <p>
              Dimulai {event.startsAt.toLocaleDateString()} {event.startsAt.toLocaleTimeString()}.
              Durasi {durationText(event.startsAt, event.endsAt)}.
            </p>
          </Card>
        </Link>
      ))}
    </div>
  );
};

const EventListPage = (props: EventListPageProps) => {
  const { createEvent, events, getEvents, user } = props;

  React.useEffect(() => { document.title = 'Exams'; }, []);
  React.useEffect(() => { getEvents(); }, [getEvents]);

  const [isShowingCreateModal, setShowingCreateModal] = React.useState(false);
  const submitNewEvent = async (data: EventFormData) => {
    return createEvent({ id: 0, ...data } as Event)
      .then(() => {
        setShowingCreateModal(false);
        getEvents();
      })
      .catch((err) => {
        if (err instanceof CharonFormError) {
          throw err.asSubmissionError();
        } else {
          throw new SubmissionError({ _error: "Unknown error" });
        }
      });
  };

  return (
    <div className="event-page">

      <Modal isShowing={isShowingCreateModal} closeModal={() => { setShowingCreateModal(false); }}>
        <h1 className="create-event-modal-title">Buat Ujian</h1>
        <EventForm onSubmit={submitNewEvent} />
      </Modal>

      <h1>Daftar Ujian</h1>

      {(user.role === USER_ROLE.ADMIN || user.role === USER_ROLE.ORGANIZER || user.role === USER_ROLE.LOCAL) && (
        <div className="create-button-row">
          <Button onClick={() => setShowingCreateModal(true)}>
            <i className="fas fa-plus"></i><span>TAMBAH</span>
          </Button>
        </div>
      )}
      {renderEvents(events)}

    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  events: charonExamSelectors.getEvents(state),
});

const mapDispatchToProps = {
  getEvents: charonExamActions.getEvents,
  createEvent: charonExamActions.createEvent,
};

export default connect(mapStateToProps, mapDispatchToProps)(EventListPage);
