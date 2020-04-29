import React from 'react';
import { Card, TextInput, Button } from 'react-hephaestus';

import { Event } from '../../../modules/charon/exam/api';
import { durationText } from '../../../modules/util/time';
import EventDetailLoading from './EventDetailLoading';
import './EventDetail.scss';
import { CharonFormError } from '../../../modules/charon/http';
import LoadingCircle from '../../loading/Circle';

interface EventDetailProps {
  event: Event | null;
  showParticipationKeyPrompt: boolean;
  onSubmitParticipatioNkey: (participationKey: string) => Promise<void>;
};

const EventDetail = (props: EventDetailProps) => {
  const { event, onSubmitParticipatioNkey, showParticipationKeyPrompt } = props;

  const [participationKey, setParticipationKey] = React.useState('');
  const [error, setError] = React.useState('');
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  if (!event) return <EventDetailLoading />;

  const handleSubmitParticipatioNkey = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);
    onSubmitParticipatioNkey(participationKey)
      .then(() => {
        setError('');
        setIsSubmitting(false);
      })
      .catch((e) => {
        setIsSubmitting(false);
        if (e instanceof CharonFormError) {
          setError('Password salah');
        }
      });
  };

  const duration = durationText(event.startsAt, event.endsAt);
  return (
    <Card className="event-detail">
      <p>{event.description}</p>
      <ul>
        <li>Waktu mulai: {event.startsAt.toLocaleDateString()} {event.startsAt.toLocaleTimeString()}</li>
        <li>Durasi: {duration}</li>
      </ul>
      {showParticipationKeyPrompt && (
        <div>
          <p>Untuk memulai ujian, silakan masukkan password ujian:</p>
          <form>
            <TextInput errorText={error} value={participationKey} onChange={(e) => setParticipationKey(e.currentTarget.value)} />
            <Button onClick={handleSubmitParticipatioNkey} type="submit">
              {isSubmitting
              ? <LoadingCircle />
              : <strong>SUBMIT</strong>}
            </Button>
          </form>
        </div>
      )}
    </Card>
  );
};

EventDetail.defaultProps = {
  event: null,
  showParticipationKeyPrompt: false,
  onSubmitParticipatioNkey: () => Promise.reject,
};

export default EventDetail;
