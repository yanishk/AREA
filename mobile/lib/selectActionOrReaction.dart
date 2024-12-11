import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import "custom_widget.dart";
import "globalData.dart";
import "spotifyAuth.dart";

class SelectActionOrReactionPage extends StatelessWidget {
  const SelectActionOrReactionPage({
    super.key,
    this.isAnAction = true,
    required this.serviceIndex,
  });
  final bool isAnAction;
  final int serviceIndex;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: SelectActionOrReactionAPICall(
          isAnAction: isAnAction, serviceIndex: serviceIndex),
    );
  }
}

class SelectActionOrReactionAPICall extends StatefulWidget {
  const SelectActionOrReactionAPICall({
    super.key,
    this.isAnAction = true,
    required this.serviceIndex,
  });
  final bool isAnAction;
  final int serviceIndex;

  @override
  SelectActionOrReaction createState() => SelectActionOrReaction(
      isAnAction: isAnAction, serviceIndex: serviceIndex);
}

class SelectActionOrReaction extends State<SelectActionOrReactionAPICall> {
  SelectActionOrReaction({
    required this.isAnAction,
    required this.serviceIndex,
  });
  final bool isAnAction;
  final int serviceIndex;
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  List<String> serviceNameList = Globaldata.actionServicesName;
  List<String> serviceColorList = Globaldata.actionServicesColor;
  List<String> serviceIconsList = Globaldata.actionServicesIcons;
  List<ActionOrReactionData> actionOrReactionList = [];
  Color colorTheme = Colors.black;
  String title = "Select action";
  String serviceName = "";
  String result = '';
  String apiRoute = "actions";
  String actionDescription = "";
  bool isInit = false;
  int actionOrReactionIndex = 0;
  List<Widget> ActionOrReactionButtonsList = [];

  void storeServiceDescription(dynamic responseData) {
    developer.log("getServices OK :\n$responseData");
    if (responseData is Map<String, dynamic> &&
        responseData[apiRoute] is List) {
      final actions = responseData[apiRoute];
      int i = 0;
      while (actions.isNotEmpty) {
        final actionResponse = actions[i];
        if (actionResponse is Map<String, dynamic>) {
          actionDescription = actionResponse['Description'] ?? 'No Description';
          actionOrReactionList.add(ActionOrReactionData.optional(
            actionResponse['Id'] ?? '',
            actionResponse['Name'] ?? '',
            actionResponse['Description'] ?? '',
            actionResponse['ServiceId'] ?? '',
            actionResponse['NbParam'] ?? 0,
          ));
          developer.log("pass");
          ActionOrReactionButtonsList.add(
            actionOrReactionButtons(actionOrReactionList.length - 1),
          );
        }
        i++;
      }
    } else {
      result = 'Invalid response structure';
    }
  }

  Future<void> getServices() async {
    developer.log('SERVER_LINK : $apiUrl');
    try {
      developer.log(
          'SERVER_LINK + redirection : ${apiUrl + serviceNameList[serviceIndex]}/$apiRoute');
      final response = await http.get(
        Uri.parse('${apiUrl + serviceNameList[serviceIndex]}/$apiRoute'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );
      final responseData = jsonDecode(response.body);
      if (response.statusCode == 200) {
        setState(() {
          storeServiceDescription(responseData);
        });
      } else {
        developer.log(
            "getServices FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
        setState(() {
          result = '${responseData['error']}';
        });
      }
    } catch (e) {
      developer.log("getServices ERROR: $e", level: 1);
      setState(() {
        result = '$e';
        developer.log(result);
        if (result.split(' ').contains("SocketException:") == true ||
            (result.split(' ').contains("Invalid") == true &&
                result.split(' ').contains("argument(s):") == true)) {
          result = Globaldata.ErrorServerConnection;
        }
      });
    }
  }

  Future<void> setReactionService() async {
    setState(() {
      serviceNameList = Globaldata.reactionServicesName;
      serviceColorList = Globaldata.reactionServicesColor;
      serviceIconsList = Globaldata.reactionServicesIcons;
      apiRoute = "reactions";
    });
  }

  Future<void> getSpotifyAuthLink() async {
    final AuthLink = await GetApi.get("authentication/Spotify");
    Globaldata.spotifyConnectionLink = AuthLink["auth-url"];
    developer.log(
        "getSpotifyAuthLink response : ${Globaldata.spotifyConnectionLink}");
    await Navigator.push(Globaldata.myContext,
        MaterialPageRoute(builder: (context) => SpotifyAuthPage()));
  }

  Future<void> saveSelectedActionOrReaction(int index) async {
    if (isAnAction) {
      Globaldata.serviceAction = actionOrReactionList[index];
      Navigator.pop(Globaldata.myContext);
    } else {
      Globaldata.serviceReaction = actionOrReactionList[index];
      Navigator.pop(Globaldata.myContext);
    }
  }

  Widget connectButton() {
    if (serviceNameList[serviceIndex] == "Spotify") {
      return CustomButton(
        backgroundColor: colorTheme,
        text: AutoSizeText(
          "Connect",
          style: TextStyle(
              fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
        ),
        onPressed: getSpotifyAuthLink,
        size: Size(double.infinity, 60),
        padding: EdgeInsets.symmetric(horizontal: 15.0),
      );
    } else {
      return SizedBox();
    }
  }

  Future<void> setParameter(int index, String name, dynamic param) async {
    actionOrReactionList[index].parameterList[name] = param;
  }

  Future<void> selectTimePopup(int index) async {
    final TimeOfDay? selectedTime = await showTimePicker(
      context: Globaldata.myContext,
      initialEntryMode: TimePickerEntryMode.input,
      initialTime: TimeOfDay.now(),
      builder: (BuildContext context, Widget? child) {
        return Theme(
          data: ThemeData.light().copyWith(
            colorScheme: ColorScheme.light(
              primary: colorTheme,
              onSurface: Colors.black,
            ),
          ),
          child: MediaQuery(
            data: MediaQuery.of(Globaldata.myContext)
                .copyWith(alwaysUse24HourFormat: true),
            child: child!,
          ),
        );
      },
    );
    if (selectedTime != null) {
      actionOrReactionList[index].parameter =
        "${selectedTime.hour.toString()}:${selectedTime.minute.toString()}";
      setParameter(index, "hour", selectedTime.hour);
      setParameter(index, "minute", selectedTime.minute);
      saveSelectedActionOrReaction(index);
    }
  }

  Widget selectDateTime(int index) {
    return CustomButton(
      backgroundColor: colorTheme,
      text: AutoSizeText(
        actionOrReactionList[index].name,
        style: TextStyle(
            fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
      ),
      onPressed: () {
        selectTimePopup(index);
      },
      size: Size(double.infinity, 60),
      padding: EdgeInsets.symmetric(horizontal: 15.0),
    );
  }

  Widget actionOrReactionButtons(int index) {
    if (actionOrReactionList[index].name != "Every day at") {
      return Column(children: [
        Padding(padding: EdgeInsets.only(top: 30)),
        CustomButton(
          backgroundColor: HexColor(serviceColorList[serviceIndex]),
          text: AutoSizeText(
            actionOrReactionList[index].name,
            style: TextStyle(
                fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),
          ),
          onPressed: () {
            saveSelectedActionOrReaction(index);
          },
          size: Size(MediaQuery.sizeOf(Globaldata.myContext).width, 60),
          padding: EdgeInsets.symmetric(horizontal: 15.0),
        )
      ]);
    }
    return selectDateTime(index);
  }

  @override
  Widget build(BuildContext context) {
    serviceName = serviceNameList[serviceIndex];
    colorTheme = HexColor(serviceColorList[serviceIndex]);
    if (!isInit) {
      if (isAnAction != true) {
        setReactionService();
      }
      isInit = true;
      getServices();
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            ActionOrReactionDescription(
                colorTheme: colorTheme,
                title: title,
                serviceIconsList: serviceIconsList,
                serviceName: serviceName,
                actionDescription: actionDescription),
            Padding(padding: EdgeInsets.only(top: 50)),
            connectButton(),
            Column(
              children: ActionOrReactionButtonsList,
            ),
          ],
        ),
      ),
    );
  }
}

class ActionOrReactionDescription extends StatelessWidget {
  const ActionOrReactionDescription({
    super.key,
    required this.colorTheme,
    required this.title,
    required this.serviceIconsList,
    required this.serviceName,
    required this.actionDescription,
  });

  final Color colorTheme;
  final String title;
  final List<String> serviceIconsList;
  final String serviceName;
  final String actionDescription;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: colorTheme,
        border: Border.all(
          color: colorTheme,
          width: 2.0,
        ),
        borderRadius: BorderRadius.circular(4.0),
      ),
      child: InkWell(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectActionOrReactionTopBar(title: title, color: Colors.white),
            Padding(padding: EdgeInsets.only(top: 20)),
            SafeWebImage(
              url:
                  "${Globaldata.domainName}${serviceIconsList[Globaldata.firstService]}",
              width: 75,
              height: 75,
            ),
            Padding(padding: EdgeInsets.only(top: 5)),
            AutoSizeText(
              serviceName,
              style: TextStyle(
                fontSize: 25,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
              maxLines: 1,
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 15.0),
              child: AutoSizeText(
                actionDescription,
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
                textAlign: TextAlign.center,
                maxLines: 10,
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
          ],
        ),
      ),
    );
  }
}

class SelectActionOrReactionTopBar extends StatelessWidget {
  const SelectActionOrReactionTopBar({
    super.key,
    required this.title,
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(
              Icons.arrow_back,
              size: 32,
              color: color,
            ),
          ),
        ),
        Center(
          child: AutoSizeText(
            title,
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: color,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}

class HexColor extends Color {
  static int _getColorFromHex(String hexColor) {
    hexColor = hexColor.toUpperCase().replaceAll("#", "");
    if (hexColor.length == 6) {
      hexColor = "FF$hexColor";
    }
    return int.parse(hexColor, radix: 16);
  }

  HexColor(final String hexColor) : super(_getColorFromHex(hexColor));
}
