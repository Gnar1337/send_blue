import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { type Client, type ClientLead, type MessageQueue, type MessageEventHistory, fetchClientData } from '../data.ts'

export const clientStore = defineStore('client',{
  state: (): {
    clients: Client[],
    currClient: Client,
    clientLeads: ClientLead[],
    allMessagesSent: MessageQueue[],
    queuedMessages: MessageQueue[],
  } => ({
    clients: [],
    currClient: {
      uid: '',
      name: '',
      messagesSent: 0,
      leads: [],
      allMessagesSent: [],
      messageQueue: []

    },   
    clientLeads: [],
    allMessagesSent: [],
    queuedMessages: []
  }),
  actions: {
    setClients(clients: Client[]) {
      this.clients = clients;
    },
    getClients(): Client[] {
      return this.clients;
    },
 async setCurrClient(client: string) {
      console.log('Setting current client to:', client);
      fetchClientData(client).then(data => {
        console.log("data got = ", data)
        this.currClient.leads = data.leads
        this.currClient.messageQueue = data.messageQueue
        this.currClient.allMessagesSent = data.allMessagesSent
        this.currClient.uid = data.uid
        this.currClient.name = data.name
      } )
    }
  },

} );
