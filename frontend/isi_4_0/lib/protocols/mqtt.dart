import 'dart:async';
import 'dart:io';
import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';

class MQTTClient {
  late MqttServerClient _client;
  int _pongCount = 0;

  Future<void> connect(String server, String clientIdentifier) async {
    _client = MqttServerClient(server, clientIdentifier);
    _client.logging(on: true);

    _client.setProtocolV311();
    _client.keepAlivePeriod = 20;
    _client.connectTimeoutPeriod = 2000;
    _client.onDisconnected = _onDisconnected;
    _client.onConnected = _onConnected;
    _client.onSubscribed = _onSubscribed;
    _client.pongCallback = _pong;

    final connMess = MqttConnectMessage()
        .withClientIdentifier(clientIdentifier)
        .startClean()
        .withWillQos(MqttQos.atLeastOnce);
    print('MQTTClient::Connecting to server $server');
    _client.connectionMessage = connMess;

    try {
      await _client.connect();
    } on NoConnectionException catch (e) {
      print('MQTTClient::Connection error - $e');
      _client.disconnect();
      rethrow;
    } on SocketException catch (e) {
      print('MQTTClient::Socket error - $e');
      _client.disconnect();
      rethrow;
    }

    if (_client.connectionStatus!.state == MqttConnectionState.connected) {
      print('MQTTClient::Connected to server $server');
    } else {
      print(
          'MQTTClient::Connection to server $server failed - status: ${_client.connectionStatus}');
      _client.disconnect();
      exit(-1);
    }
  }

  void subscribe(String topic, MqttQos qos) {
    _client.subscribe(topic, qos);
  }

  void unsubscribe(String topic) {
    _client.unsubscribe(topic);
  }

  Future<void> sleep(int seconds) async {
    await MqttUtilities.asyncSleep(seconds);
  }

  void publish(String topic, MqttQos qos, String message) {
    final builder = MqttClientPayloadBuilder();
    builder.addString(message);
    _client.publishMessage(topic, qos, builder.payload!);
  }

  void disconnect() {
    _client.disconnect();
  }

  void _onSubscribed(String topic) {
    print('MQTTClient::Subscription confirmed for topic $topic');
  }

  void _onDisconnected() {
    print('MQTTClient::Disconnected from server');
    if (_client.connectionStatus!.disconnectionOrigin ==
        MqttDisconnectionOrigin.solicited) {
      print('MQTTClient::Disconnection was solicited');
    } else {
      print('MQTTClient::Disconnection was unsolicited');
    }
  }

  void _onConnected() {
    print('MQTTClient::Connected to server');
  }

  void _pong() {
    _pongCount++;
    print('MQTTClient::Pong received ($_pongCount)');
  }
}
