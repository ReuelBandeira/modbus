class DataProtocols {
  int id;
  String protocol;
  String alias;
  int version;
  Map<String, dynamic> config;

  DataProtocols(
      {required this.id,
      required this.protocol,
      required this.alias,
      required this.version,
      required this.config});
}
