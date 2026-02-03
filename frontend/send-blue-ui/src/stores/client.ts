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
    getCurrClient(): Client {
      return this.currClient;
    },
 async setCurrClient(client: string) {
      fetchClientData(client).then(data => {
        this.currClient.leads = data.leads
        this.currClient.messageQueue = data.messageQueue
        this.currClient.allMessagesSent = data.allMessagesSent
        this.currClient.uid = data.uid
        this.currClient.name = data.name
        this.currClient.messagesSent = data.allMessagesSent.length
      } )
    }
  },

} );
