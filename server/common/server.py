import socket
import logging
import signal
from common.utils import receive_bets, store_bets, load_bets, has_won, write_to_socket

class Server:
    def __init__(self, port, listen_backlog, agencies_amount):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._clients = []
        self.agencies = {}
        self.agencies_amount = agencies_amount

        signal.signal(signal.SIGTERM, self.__handle_shutdown)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self._clients.append(client_sock)
                self.__handle_client_connection(client_sock)
            except OSError as e:
                if not self._running:
                    logging.error(f"action: accept_connections | result: fail | error: {e}")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bets = receive_bets(client_sock)
            store_bets(bets)
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        except OSError as e:
            logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
        finally:
            self.__send_bet_results(client_sock, bets[0].agency)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
    
    def __handle_shutdown(self, _signum, _frame):
        self._running = False
        for client in self._clients:
            client.close()
        self._server_socket.close()
        logging.info("action: shutdown | result: success")

    def __send_bet_results(self, client_sock, agency_id):
        self.agencies[agency_id] = client_sock
        if len(self.agencies) != self.agencies_amount:
            return
        try:
            bets = load_bets()
            agency_winers = {}
            for bet in bets:
                if has_won(bet):
                    agency_winers[bet.agency] = agency_winers.get(bet.agency, []) + [bet.document]
            for agency_id in self.agencies.keys():
                winers = agency_winers.get(agency_id, [])
                documents_str = ';'.join(winers) + '\n' if winers else '\n'
                write_to_socket(self.agencies[agency_id], documents_str.encode())
            logging.info("action: sorteo | result: success")
        except OSError as e:
            logging.error(f'action: sorteo | result: fail | error: {e}')
        finally:
            for client in self.agencies.values():
                client.close()
            self.agencies = {}