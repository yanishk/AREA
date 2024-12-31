import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../custom_widget.dart";
import "../globalData.dart";
import 'createWidget.dart';
import 'createWorkflowsData.dart';
import "../auth/connectSpotify.dart";

class SelectActionOrReactionPage extends StatelessWidget {
  const SelectActionOrReactionPage({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: SelectActionOrReactionAPICall(serviceIndex: serviceIndex),
    );
  }
}

class SelectActionOrReactionAPICall extends StatefulWidget {
  const SelectActionOrReactionAPICall({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  SelectActionOrReaction createState() => SelectActionOrReaction(serviceIndex: serviceIndex);
}

class SelectActionOrReaction extends State<SelectActionOrReactionAPICall> {
  SelectActionOrReaction({
    required this.serviceIndex,
  });
  final int serviceIndex;
  List<ActionOrReactionData> actionOrReactionList = [];
  Color colorTheme = Colors.black;
  String apiRoute = "actions";
  String actionDescription = "";
  List<Widget> ActionOrReactionButtonsList = [];
  Size buttonSize = Size(MediaQuery.sizeOf(Globaldata.myContext).width, 20);
  double buttonRadius = 8;

  @override
  void initState() {
    super.initState();
    if (!isAnAction) {
      setState(() {
        apiRoute = "reactions";
      });
    }
    getActionsOrReactions();
  }

  Future<void> getActionsOrReactions() async {

    setState(() {
      if (apiRoute == "actions") {
        storeServiceDescription(Globaldata.actionServicesActions);
      } else {
        storeServiceDescription(Globaldata.reactionServicesReaction);
      }
    });
  }

  void storeServiceDescription(List<ActionOrReactionList> allActionOrReactionList) {

    for (int i = 0; i < allActionOrReactionList.length; i++) {
      if (allActionOrReactionList[i].servicename == CreateData.serviceNameList[serviceIndex]) {
        actionOrReactionList.add(ActionOrReactionData.optional(
          allActionOrReactionList[i].id,
          allActionOrReactionList[i].name,
          allActionOrReactionList[i].description,
          allActionOrReactionList[i].serviceid,
          allActionOrReactionList[i].nbparam,
        ));
        ActionOrReactionButtonsList.add(
          actionOrReactionButtons(actionOrReactionList.length - 1),
        );
      }
    }
  }

  Widget actionOrReactionButtons(int index) {
    if (actionOrReactionList[index].name == "Every day at") {
      return selectDateTime(index);
    }
    return Column(children: [
      Padding(padding: EdgeInsets.only(top: 15)),
      CustomButton(
        backgroundColor: HexColor(CreateData.serviceColorList[serviceIndex]),
        text: ButtonText(text: actionOrReactionList[index].name,),
        onPressed: () {
          saveSelectedActionOrReaction(index);
        },
        size: buttonSize,
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        borderCircularRadius: buttonRadius,
      )
    ]);
  }

  Widget selectDateTime(int index) {
    return Padding(
      padding: EdgeInsets.only(top: 15),
      child: CustomButton(
        backgroundColor: HexColor(CreateData.serviceColorList[serviceIndex]),
        text: ButtonText(text: actionOrReactionList[index].name,),
        onPressed: () {
          selectTimePopup(index);
        },
        size: buttonSize,
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        borderCircularRadius: buttonRadius,
      ),
    );
  }

  Future<void> saveSelectedActionOrReaction(int index) async {
    if (isAnAction) {
      CreateData.serviceAction = actionOrReactionList[index];
    } else {
      CreateData.serviceReaction = actionOrReactionList[index];
    }
    Navigator.pop(Globaldata.myContext);
  }

  Widget connectButton() {
    if (CreateData.serviceNameList[serviceIndex] == "Spotify") {
      return CustomButton(
        backgroundColor: colorTheme,
        text: ButtonText(text: "Connect",),
        onPressed: () {MySpotifyAuth.connectSpotify("service");},
        size: buttonSize,
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        borderCircularRadius: buttonRadius,
      );
    } else {
      return SizedBox();
    }
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
      actionOrReactionList[index].parameterList["hour"] = selectedTime.hour;
      actionOrReactionList[index].parameterList["minute"] = selectedTime.minute;
      saveSelectedActionOrReaction(index);
    }
  }

  @override
  Widget build(BuildContext context) {
    colorTheme = HexColor(CreateData.serviceColorList[serviceIndex]);
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            ActionOrReactionDescription(
                colorTheme: colorTheme,
                title: "Select $apiRoute",
                serviceIconsList: CreateData.serviceIconsList,
                serviceName: CreateData.serviceNameList[serviceIndex],
                actionDescription: actionDescription),
            Padding(padding: EdgeInsets.only(top: 15)),
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

class ButtonText extends StatelessWidget {
  const ButtonText({
    super.key,
    required this.text,
  });

  final String text;

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: Alignment.centerLeft,
      child: new AutoSizeText(
        text,
        style: TextStyle(
          fontSize: 20,
          fontWeight: FontWeight.bold,
          color: Colors.white,
          height: 2.3,
        ),
      )
    );
  }
}
