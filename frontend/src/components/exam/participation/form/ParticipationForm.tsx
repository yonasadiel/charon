import parseInt from 'lodash/parseInt';
import React from 'react';
import { Field, reduxForm, InjectedFormProps, WrappedFieldProps } from 'redux-form';
import { Button, BUTTON_TYPE_CONTAINED, TextInput } from 'react-hephaestus';

import { User } from '../../../../modules/charon/auth/api';
import { Participation, Venue } from '../../../../modules/charon/exam/api';
import LoadingCircle from '../../../loading/Circle';
import './ParticipationForm.scss';

export interface ParticipationFormData extends Partial<Participation> {};
export interface ParticipationFormProps {
  users: User[];
  venues: Venue[];
}

const VenueField = (props: WrappedFieldProps & { venues: Venue[] }) => {
  const { input, meta: { error }, venues } = props;
  return (
    <div className="venue-select">
      <select {...input}>
        <option></option>
        {venues.map((venue) => (
          <option value={venue.id} key={venue.id}>{venue.name}</option>
        ))}
      </select>
      <small className="error">{error}</small>
    </div>
  );
};

const UserField = (props: WrappedFieldProps & { users: User[] }) => {
  const { input, meta: { error }, users } = props;
  return (
    <div className="user-select">
      <select {...input}>
        <option></option>
        {users.map((user) => (
          <option value={user.username} key={user.username}>{user.name}</option>
        ))}
      </select>
      <small className="error">{error}</small>
    </div>
  );
};

const renderKeyField = ({ input, meta: { error } }: WrappedFieldProps) => (
  <div className="key">
    <TextInput {...input} />
    <small className="error">{error}</small>
  </div>
)

const ParticipationForm = (props: ParticipationFormProps & InjectedFormProps<ParticipationFormData, ParticipationFormProps>) => {
  const { error, handleSubmit, reset, submitting, users, venues } = props;
  return (
    <form className="participation-form" onSubmit={(e) => { handleSubmit(e); reset() }}>
      <div className="field-row">
        <Field name="userUsername" users={users} component={UserField} />
        <Field name="venueId" venues={venues} component={VenueField} parse={parseInt} />
        <Field name="key" component={renderKeyField} />
        <Button buttonType={BUTTON_TYPE_CONTAINED} type="submit" className="button">
          { !submitting ? (<strong>BUAT</strong>) : <LoadingCircle /> }
        </Button>
      </div>
      <small className="error">{ error }</small>
    </form>
  );
};

export default reduxForm<ParticipationFormData, ParticipationFormProps>({ form: 'participation' })(ParticipationForm);
