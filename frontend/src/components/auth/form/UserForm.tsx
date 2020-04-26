import React from 'react';
import { Field, reduxForm, InjectedFormProps, WrappedFieldProps } from 'redux-form';
import { Button, TextInput, BUTTON_TYPE_CONTAINED } from 'react-hephaestus';

import { User } from '../../../modules/charon/auth/api';
import LoadingCircle from '../../loading/Circle';
import './UserForm.scss';

export interface UserFormData extends Partial<User> {
};

export interface UserFormProps extends InjectedFormProps<UserFormData> {
};

const renderNameField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Nama" type="text" />
);

const renderUsernameField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Username" type="text" />
);

const renderPasswordField = (field: WrappedFieldProps) => (
  <TextInput {...field.input} errorText={field.meta.error} placeholder="Password" type="password" />
);

const renderRoleField = (field: WrappedFieldProps) => (
  <div className="role-select-row">
    <select {...field.input}>
      <option value="organizer">Penyelenggara</option>
      <option value="local">Panitia lokal</option>
      <option value="participant">Peserta ujian</option>
    </select>
    <p>{field.meta.error}</p>
  </div>
);

const UserForm = (props: UserFormProps) => {
  const { error, handleSubmit, submitting } = props;
  return (
    <form className="user-form" onSubmit={handleSubmit}>
      <Field name="name" component={renderNameField} />
      <Field name="username" component={renderUsernameField} />
      <Field name="password" component={renderPasswordField} />
      <Field name="role" component={renderRoleField} />

      <small className="error">{ error }</small>

      <div className="user-button-row">
        <Button buttonType={BUTTON_TYPE_CONTAINED} type="submit">
          { !submitting ? (<strong>BUAT</strong>) : <LoadingCircle /> }
        </Button>
      </div>
    </form>
  );
};

export default reduxForm<UserFormData>({ form: 'user' })(UserForm);
