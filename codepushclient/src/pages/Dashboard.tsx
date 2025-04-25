import React from 'react';
import {
    Container,
    Box,
    Typography,
    Paper,
    Button,
    Grid,
    Card,
    CardContent,
    Divider,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';

const Dashboard: React.FC = () => {
    const { user, logout } = useAuth();

    if (!user) {
        return null;
    }

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
    };

    return (
        <Container maxWidth="lg">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <Paper sx={{ p: 3 }}>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="h4" component="h1">
                                    Welcome, {user.username}!
                                </Typography>
                                <Button variant="outlined" color="error" onClick={logout}>
                                    Logout
                                </Button>
                            </Box>
                        </Paper>
                    </Grid>

                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Your App ID
                                </Typography>
                                <Typography variant="body1" sx={{ wordBreak: 'break-all' }}>
                                    {user.appId}
                                </Typography>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    sx={{ mt: 2 }}
                                    onClick={() => copyToClipboard(user.appId)}
                                >
                                    Copy App ID
                                </Button>
                            </CardContent>
                        </Card>
                    </Grid>

                    <Grid item xs={12} md={6}>
                        <Card>
                            <CardContent>
                                <Typography variant="h6" gutterBottom>
                                    Your Token
                                </Typography>
                                <Typography variant="body1" sx={{ wordBreak: 'break-all' }}>
                                    {user.token}
                                </Typography>
                                <Button
                                    variant="outlined"
                                    size="small"
                                    sx={{ mt: 2 }}
                                    onClick={() => copyToClipboard(user.token)}
                                >
                                    Copy Token
                                </Button>
                            </CardContent>
                        </Card>
                    </Grid>

                    <Grid item xs={12}>
                        <Paper sx={{ p: 3 }}>
                            <Typography variant="h6" gutterBottom>
                                Account Information
                            </Typography>
                            <Divider sx={{ my: 2 }} />
                            <Grid container spacing={2}>
                                <Grid item xs={12} sm={6}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                        Email
                                    </Typography>
                                    <Typography variant="body1">{user.email}</Typography>
                                </Grid>
                                {user.companyName && (
                                    <Grid item xs={12} sm={6}>
                                        <Typography variant="subtitle2" color="text.secondary">
                                            Company Name
                                        </Typography>
                                        <Typography variant="body1">{user.companyName}</Typography>
                                    </Grid>
                                )}
                                {user.phoneNumber && (
                                    <Grid item xs={12} sm={6}>
                                        <Typography variant="subtitle2" color="text.secondary">
                                            Phone Number
                                        </Typography>
                                        <Typography variant="body1">{user.phoneNumber}</Typography>
                                    </Grid>
                                )}
                                <Grid item xs={12} sm={6}>
                                    <Typography variant="subtitle2" color="text.secondary">
                                        Account Created
                                    </Typography>
                                    <Typography variant="body1">
                                        {new Date(user.createdAt).toLocaleDateString()}
                                    </Typography>
                                </Grid>
                            </Grid>
                        </Paper>
                    </Grid>
                </Grid>
            </Box>
        </Container>
    );
};

export default Dashboard; 