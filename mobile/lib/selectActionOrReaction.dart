import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import "package:http/src/response.dart";
import 'dart:convert';

import "custom_widget.dart";
import "globalData.dart";
import "spotifyAuth.dart";
import "apiCall/apiRequest.dart";

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
    developer.log("storeServiceDescription, service '$apiRoute' :\n$responseData");
    if (responseData is Map<String, dynamic> &&
        responseData[apiRoute] is List) {
      List<dynamic> actions = responseData[apiRoute];
      while (!actions.isEmpty) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          actionDescription = actionResponse['Description'] ?? 'No Description';
          actionOrReactionList.add(ActionOrReactionData.optional(
            actionResponse['Id'] ?? '',
            actionResponse['Name'] ?? '',
            actionResponse['Description'] ?? '',
            actionResponse['ServiceId'] ?? '',
            actionResponse['NbParam'] ?? 0,
          ));
          ActionOrReactionButtonsList.add(
            actionOrReactionButtons(actionOrReactionList.length - 1),
          );
        }
      }
    } else {
      result = 'Invalid response structure';
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
    Response response = await ApiRequest.get("authentication/Spotify");
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getSpotifyAuthLink FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    Globaldata.spotifyConnectionLink = responseData["auth-url"];
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

  Future<void> getActionsOrReactions() async {
    Response response = await ApiRequest.get("${serviceNameList[serviceIndex]}/$apiRoute");
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getSpotifyAuthLink FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    setState(() {
      storeServiceDescription(responseData);
    });
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
      getActionsOrReactions();
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
