import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../../custom_widget.dart";
import "../../globalData.dart";

class DisplayActionOrReactionPage extends StatelessWidget {
  const DisplayActionOrReactionPage({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: DisplayActionOrReactionAPICall(serviceIndex: serviceIndex),
    );
  }
}

class DisplayActionOrReactionAPICall extends StatefulWidget {
  const DisplayActionOrReactionAPICall({
    super.key,
    required this.serviceIndex,
  });
  final int serviceIndex;

  @override
  DisplayActionOrReaction createState() => DisplayActionOrReaction(serviceIndex: serviceIndex);
}

class DisplayActionOrReaction extends State<DisplayActionOrReactionAPICall> {
  DisplayActionOrReaction({
    required this.serviceIndex,
  });
  final int serviceIndex;
  Color colorTheme = Colors.black;
  String actionDescription = "";
  List<Widget> ActionOrReactionButtonsList = [];
  double buttonRadius = 8;
  List<ActionOrReactionList> actionList = [];
  List<ActionOrReactionList> reactionList = [];

  @override
  void initState() {
    super.initState();
    actionList = [];
    reactionList = [];
    storeActionOrReactionsServices(Globaldata.actionServicesActions, actionList);
    storeActionOrReactionsServices(Globaldata.reactionServicesReaction, reactionList);
    fillActionOrReactionButtons("Actions", actionList);
    if (actionList.isNotEmpty) {
      ActionOrReactionButtonsList.add(
        Padding(padding: EdgeInsets.only(top: 30))
      );
    }
    fillActionOrReactionButtons("Reactions", reactionList);
    actionDescription = Globaldata.serviceList[serviceIndex].description;
  }

  void storeActionOrReactionsServices(List<ActionOrReactionList> src, List<ActionOrReactionList> dest) {
    for (int i = 0; i < src.length; i++) {
      if (src[i].servicename == Globaldata.serviceList[serviceIndex].name) {
        dest.add(src[i]);
      }
    }
  }

  void fillActionOrReactionButtons(String name, List<ActionOrReactionList> actionOrReactionList) {
    if (actionOrReactionList.isNotEmpty) {
      ActionOrReactionButtonsList.add(
        LeftWidget(
          child: AutoSizeText(
            name,
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.w900,
              color: Colors.black,
            ),
            maxLines: 1,
          ),
        ),
      );
      for (int i = 0; i < actionOrReactionList.length; i++) {
        ActionOrReactionButtonsList.add(
          actionOrReactionButtons(i, actionOrReactionList)
        );
      }
    }
  }

  Widget actionOrReactionButtons(int index, List<ActionOrReactionList> actionOrReactionList) {
    return Column(children: [
      Padding(padding: EdgeInsets.only(top: 15)),
      Padding(
        padding: EdgeInsets.symmetric(horizontal: 15.0),
        child: TextWithBackground(
          text: actionOrReactionList[index].name,
          backgroundColor: HexColor(Globaldata.serviceList[serviceIndex].color),
          textColor: Colors.white,
          fontSize: 20,
          borderRadius: buttonRadius,
          height: 45,
          width: MediaQuery.sizeOf(Globaldata.myContext).width,
          alignment: Alignment.centerLeft,
        ),
      )
    ]);
  }

  @override
  Widget build(BuildContext context) {
    colorTheme = HexColor(Globaldata.serviceList[serviceIndex].color);
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            ActionOrReactionDescription(
                colorTheme: colorTheme,
                serviceIcon: Globaldata.serviceList[serviceIndex].icon,
                serviceName: Globaldata.serviceList[serviceIndex].name,
                actionDescription: actionDescription),
            Padding(padding: EdgeInsets.only(top: 15)),
            AutoSizeText(
              "Details",
              style: TextStyle(
                fontSize: 25,
                fontWeight: FontWeight.w900,
                color: Colors.black,
              ),
              maxLines: 1,
            ),
            Column(
              children: ActionOrReactionButtonsList,
            ),
            Padding(padding: EdgeInsets.only(top: 15))
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

class TopBar extends StatelessWidget {
  const TopBar({
    super.key,
    this.color = Colors.black,
  });

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
        Center(child: AutoSizeText("",),),
      ],
    );
  }
}


class ActionOrReactionDescription extends StatelessWidget {
  const ActionOrReactionDescription({
    super.key,
    required this.colorTheme,
    required this.serviceIcon,
    required this.serviceName,
    required this.actionDescription,
    this.title = "",
  });

  final Color colorTheme;
  final String title;
  final String serviceIcon;
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
            Padding(padding: EdgeInsets.only(top: 35)),
            TopBar(color: Colors.white),
            SafeWebImage(
              url:
                  "${Globaldata.domainName}$serviceIcon",
              width: 125,
              height: 125,
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
