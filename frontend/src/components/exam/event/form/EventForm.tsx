import React from 'react';
import { Field, reduxForm, InjectedFormProps, WrappedFieldProps } from 'redux-form';
import { Button, TextInput, BUTTON_TYPE_CONTAINED } from 'react-hephaestus';

import { Event } from '../../../../modules/charon/exam/api';
import LoadingCircle from '../../../loading/Circle';
import './EventForm.scss';

export interface EventFormData extends Partial<Event> {
};

export interface EventFormProps extends InjectedFormProps<EventFormData> {
  event?: Event;
};

const renderTitleField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Judul" type="text" />
);

const renderStartsAtField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Waktu mulai" type="text" />
);

const renderEndsAtField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Waktu selesai" type="text" />
);

const EventForm = (props: EventFormProps) => {
  const { error, event, handleSubmit, submitting } = props;
  return (
    <form className="event-form" onSubmit={handleSubmit}>
      <h2 className="title">{!!event ? 'Ubah Ujian' : 'Buat Ujian'}</h2>
      <Field name="title" component={renderTitleField} />
      <Field name="startsAt" component={renderStartsAtField} />
      <Field name="endsAt" component={renderEndsAtField} />

      <small className="error">{ error }</small>

      <div className="event-button-row">
        <Button buttonType={BUTTON_TYPE_CONTAINED}>
          { !submitting ? (<strong>BUAT</strong>) : <LoadingCircle /> }
        </Button>
      </div>
    </form>
  );
};

export default reduxForm<EventFormData>({ form: 'event' })(EventForm);