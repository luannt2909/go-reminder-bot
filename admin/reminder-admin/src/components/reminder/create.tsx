import React, {cloneElement} from 'react'
import {BooleanInput, Create, SelectInput, SimpleForm, TextInput, TopToolbar} from 'react-admin'
import {WebhookTypes} from './webhookType'
import {Box, Typography} from "@material-ui/core";

const ReminderCreate = (props) => {
    return (
        <Create title='Create a Reminder' {...props}
        actions={<ListActions {...props}/>}>
            <SimpleForm sx={{ maxWidth: 1000 }}>
                <TextInput source='name' required fullWidth/>
                <BooleanInput source='is_active'/>
                <TextInput source='schedule' required />
                <SelectInput choices={WebhookTypes} required source='webhook_type'/>
                <TextInput multiline fullWidth required source='webhook'/>
                <TextInput multiline fullWidth required source='message'/>
            </SimpleForm>
        </Create>
    )
}

export default ReminderCreate