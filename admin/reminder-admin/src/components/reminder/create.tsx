import React from 'react'
import {BooleanInput, Create, SelectInput, SimpleForm, TextInput} from 'react-admin'
import {WebhookTypes} from './webhookType'
import {Box, Grid, Typography} from "@material-ui/core";

const ReminderCreate = (props) => {
    return (
        <Create title='Create a Reminder' {...props}>
            <SimpleForm sx={{maxWidth: 1000}}>
                <Typography variant="h6" gutterBottom>General</Typography>
                <TextInput source='name'
                           required
                           fullWidth
                           helperText="ex: Daily reminder bot"/>

                <BooleanInput source='is_active'/>
                <Separator/>

                <Typography variant="h6" gutterBottom>Cron Schedule Specification</Typography>
                <TextInput source='schedule'
                           sx={{width: "50%"}}
                           required
                           helperText="ex: '* * * * *', '@every 5m',... "/>
                <Separator/>
                <Typography variant="h6" gutterBottom>Webhook infomation</Typography>

                <Grid container >
                    <Grid item xs={2} >
                        <Box>
                            <SelectInput choices={WebhookTypes}
                                         required
                                         label="Webhook Type"
                                         source='webhook_type'/>
                        </Box>
                    </Grid>
                    <Grid item xs={10} >
                        <Box>
                            <TextInput sx={{pl: 1}} multiline fullWidth required source='webhook' label="Webhook URL"/>
                        </Box>
                    </Grid>
                </Grid>
                <TextInput multiline
                           fullWidth
                           required
                           placeholder="Input your message..."
                           source='message'/>
            </SimpleForm>
        </Create>
    )
}
const Separator = () => <Box pt="1em"/>;

export default ReminderCreate